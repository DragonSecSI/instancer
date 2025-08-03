from flask import current_app, Blueprint, redirect, url_for, session, request
from CTFd.plugins import register_plugin_assets_directory, register_user_page_menu_bar
from CTFd.models import db, Teams, Users
from CTFd.utils import get_config
from CTFd.utils.user import authed, get_current_user, get_current_team
from CTFd.plugins.flags import BaseFlag, FLAG_CLASSES
from sqlalchemy import ForeignKey
from sqlalchemy.types import Integer, String, Boolean, DateTime
from sqlalchemy.orm import relationship
from sqlalchemy.event import listens_for
import requests
import datetime

INSTANCER_API_URL = "https://instancer.vuln.si"
INSTANCER_ADMIN_TOKEN = "admin"

class InstancerTokenTable(db.Model):
    __tablename__ = "instancer_tokens"
    id = db.Column(Integer, primary_key=True)
    team_id = db.Column(Integer, ForeignKey("teams.id"), unique=True, nullable=True)
    user_id = db.Column(Integer, ForeignKey("users.id"), unique=True, nullable=True)
    token = db.Column(String(128), nullable=False)

    team = relationship("Teams", backref="team_instancer_token")
    user = relationship("Users", backref="user_instancer_token")

class InstancedFlagAudit(db.Model):
    __tablename__ = "instanced_flag_audit"
    id = db.Column(Integer, primary_key=True)
    team_id = db.Column(Integer, ForeignKey("teams.id"), nullable=True)
    user_id = db.Column(Integer, ForeignKey("users.id"), nullable=True)
    flag_submitted = db.Column(String(512), nullable=False)
    correct = db.Column(Boolean, nullable=False)
    active_instance = db.Column(Boolean, nullable=True)
    wrong_team = db.Column(Boolean, nullable=True)
    timestamp = db.Column(DateTime, nullable=False, default=datetime.datetime.utcnow)

class InstancedFlag(BaseFlag):
    name = "instanced"
    templates = {
        "create": "/plugins/ctfd-instancer/assets/flag/create.html",
        "update": "/plugins/ctfd-instancer/assets/flag/edit.html",
    }

    @staticmethod
    def compare(flag_obj, provided):
        user_mode = get_config("user_mode")
        user_id = None
        team_id = None
        remote_id = None
        if user_mode:
            user = get_current_user()
            if user:
                remote_id = user.id
        else:
            team = get_current_team()
            if team:
                remote_id = team.id

        api_base = INSTANCER_API_URL
        api_url = f"{api_base}/api/v1/flag/submit"  # Change as needed
        api_auth = {"Authorization": INSTANCER_ADMIN_TOKEN}
        payload = {
            "flag": provided,
            "remote_id": str(remote_id),
            "challenge_id": str(flag_obj.challenge.id),
        }

        correct = False
        active_instance = None
        wrong_team = None
        try:
            r = requests.post(api_url, json=payload, headers=api_auth, timeout=5)
            if r.status_code == 200:
                data = r.json()
                correct = data["correct"]
                active_instance = data["active_instance"]
                wrong_team = data["wrong_team"]
            else:
                current_app.logger.warning(f"Flag verification API returned {r.status_code}: {r.text}")
        except Exception as e:
            current_app.logger.error(f"Error verifying instanced flag via API: {e}")
            return False

        # Audit every submission, except exceptions
        with current_app.app_context():
            audit = InstancedFlagAudit(
                team_id=team_id,
                user_id=user_id,
                flag_submitted=provided,
                correct=correct,
                active_instance=active_instance,
                wrong_team=wrong_team,
            )
            db.session.add(audit)
        #     db.session.commit()

        return correct

FLAG_CLASSES["instanced"] = InstancedFlag

instancer_bp = Blueprint("instancer_plugin", __name__)


@instancer_bp.route("/instancer")
def instancer_portal_redirect():
    # Not logged in? Redirect to login with next param.
    if not authed():
        return redirect(url_for("auth.login", next=request.path))

    user = get_current_user()
    mode = get_config("user_mode")

    token = None
    if mode:  # User mode
        record = InstancerTokenTable.query.filter_by(user_id=user.id).first()
        if record:
            token = record.token
        else:
            token = generate_token(user.name, user.id)
            if token:
                itt = InstancerTokenTable(team_id=None, user_id=user.id, token=token)
                db.session.add(itt)
                db.session.commit()
    else:  # Team mode
        team = get_current_team()
        if team:
            record = InstancerTokenTable.query.filter_by(team_id=team.id).first()
            if record:
                token = record.token
            else:
                token = generate_token(team.name, team.id)
                if token:
                    itt = InstancerTokenTable(team_id=team.id, user_id=None, token=token)
                    db.session.add(itt)
                    db.session.commit()

    if token:
        instancer_base = INSTANCER_API_URL
        instancer_url = f"{instancer_base}?token={token}"
        return redirect(instancer_url)
    else:
        # No token found, fallback or error
        return "No token found for your account.", 400

def generate_token(name, remote_id):
    try:
        team_data = {
            "name": name,
            "remote_id": str(remote_id),
        }

        api_base = INSTANCER_API_URL
        api_url = f"{api_base}/api/v1/auth/team/register"  # Change as needed
        api_auth = {"Authorization": INSTANCER_ADMIN_TOKEN}
        response = requests.post(api_url, json=team_data, headers=api_auth, timeout=5)
        if response.status_code == 201:
            data = response.json()
            token = data.get("token")
            if token:
                with current_app.app_context():
                    if get_config("user_mode"):
                        new_token = InstancerTokenTable(team_id=None, user_id=remote_id, token=token)
                    else:
                        new_token = InstancerTokenTable(team_id=remote_id, user_id=None, token=token)
                    db.session.add(new_token)
                    current_app.logger.info(f"Saved team token for team {remote_id}")
                return token
            else:
                current_app.logger.warning(f"No token in API response: {data}")
        else:
            current_app.logger.warning(f"Instancer API returned {response.status_code}: {response.text}")
    except Exception as e:
        current_app.logger.error(f"Failed to notify instancer API: {str(e)}")

    return None

def load(app):
    register_plugin_assets_directory(app, base_path='/plugins/ctfd-instancer/assets/')
    app.register_blueprint(instancer_bp)
    register_user_page_menu_bar("Instancer", "/instancer")
    app.logger.info("CTFd Instancer plugin loaded.")

    with app.app_context():
        db.create_all()
