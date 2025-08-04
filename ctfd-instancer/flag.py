from flask import current_app
from CTFd.models import db
from CTFd.utils import get_config
from CTFd.utils.user import get_current_user, get_current_team
from CTFd.plugins.flags import BaseFlag
from sqlalchemy import ForeignKey
from sqlalchemy.types import Integer, String, Boolean, DateTime
import requests
import datetime

from .config import get_instancer_api_url, get_instancer_api_token


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
                user_id = user.id
                remote_id = user.id
        else:
            team = get_current_team()
            if team:
                team_id = team.id
                remote_id = team.id

        api_base = get_instancer_api_url()
        api_url = f"{api_base}/api/v1/flag/submit"  # Change as needed
        api_auth = {"Authorization": get_instancer_api_token()}
        payload = {
            "flag": provided,
            "remote_id": str(remote_id),
            "remote_challenge_id": str(flag_obj.challenge.id),
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

        audit = InstancedFlagAudit(
            team_id=team_id,
            user_id=user_id,
            flag_submitted=provided,
            correct=correct,
            active_instance=active_instance,
            wrong_team=wrong_team,
        )
        db.session.add(audit)
        db.session.commit()

        return correct
