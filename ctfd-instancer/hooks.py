from flask import current_app, Blueprint, redirect, url_for, request
from CTFd.models import db
from CTFd.utils import get_config
from CTFd.utils.user import authed, get_current_user, get_current_team
from sqlalchemy import ForeignKey
from sqlalchemy.types import Integer, String
from sqlalchemy.orm import relationship
import requests
import time

from .challenge import CTFdInstancedChallenge
from .config import get_instancer_api_url, get_instancer_api_token, get_instancer_public_url


def ctf_has_started():
    start = get_config("start")
    if start:
        try:
            start = int(start)
        except ValueError:
            return True

        now = int(time.time())
        return now >= start
    return True

def is_admin():
    user = get_current_user()
    if user:
        return user.type == "admin"
    return False


class InstancerTokenTable(db.Model):
    __tablename__ = "instancer_tokens"
    id = db.Column(Integer, primary_key=True)
    team_id = db.Column(Integer, ForeignKey("teams.id"), unique=True, nullable=True)
    user_id = db.Column(Integer, ForeignKey("users.id"), unique=True, nullable=True)
    token = db.Column(String(128), nullable=False)

    team = relationship("Teams", backref="team_instancer_token")
    user = relationship("Users", backref="user_instancer_token")


instancer_bp = Blueprint("instancer_plugin", __name__)

@instancer_bp.route("/instancer")
def instancer_portal_redirect():
    if not authed():
        return redirect(url_for("auth.login", next=request.path))

    if not ctf_has_started() and not is_admin():
        return redirect(url_for("challenges.listing"))

    cid = request.args.get("chall")
    if cid:
        try:
            cid = int(cid)
        except ValueError:
            return "Invalid challenge ID.", 400

        chall = CTFdInstancedChallenge.query.filter_by(id=cid).first()
        if not chall:
            return "Challenge not found or not an instanced challenge.", 404

        cid = chall.instancer_id

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
        instancer_base = get_instancer_public_url()
        instancer_url = f"{instancer_base}?token={token}"
        if cid:
            instancer_url += f"&chall={cid}"
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

        api_base = get_instancer_api_url()
        api_url = f"{api_base}/api/v1/auth/team/register"  # Change as needed
        api_auth = {"Authorization": get_instancer_api_token()}
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
