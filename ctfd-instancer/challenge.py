from flask import Blueprint
from CTFd.models import Challenges
from CTFd.plugins.challenges import CTFdStandardChallenge
from CTFd.models import Challenges, db
from CTFd.exceptions.challenges import ChallengeUpdateException, ChallengeCreateException
import requests

instanced_challenge_bp = Blueprint("instanced_challenge", __name__, template_folder="templates")
from .config import get_instancer_api_url, get_instancer_api_token

class CTFdInstancedChallenge(Challenges):
    __tablename__ = "instanced_challenges"
    __mapper_args__ = {"polymorphic_identity": "instanced"}
    id = db.Column(db.Integer, db.ForeignKey("challenges.id", ondelete="CASCADE"), primary_key=True)
    instancer_id = db.Column(db.Integer, nullable=True, unique=True)  # Unique ID for the challenge in Instancer
    challenge_type = db.Column(db.String(32), nullable=False)
    flag_base = db.Column(db.String(128), nullable=False)
    flag_type = db.Column(db.Integer, nullable=False, default=0) # Mask: 1 suffix, 2 leetify, 4 capitalize
    duration = db.Column(db.Integer, nullable=True, default=1800)  # Duration in seconds
    repository = db.Column(db.String(256), nullable=True)  # Repository URL for the challenges
    chart = db.Column(db.String(128), nullable=True)  # Chart suffix for the challenge
    chart_version = db.Column(db.String(32), nullable=True)  # Version of the chart
    values = db.Column(db.String(1024), nullable=True)  # Helm overrides in cli notation

    def __init__(self, *args, **kwargs):
        super(CTFdInstancedChallenge, self).__init__(**kwargs)

class CTFdInstancedTypeChallenge(CTFdStandardChallenge):
    id = "instanced"
    name = "instanced"
    templates = {
        "create": "/plugins/ctfd-instancer/assets/chall/create.html",
        "update": "/plugins/ctfd-instancer/assets/chall/update.html",
        "view": "/plugins/ctfd-instancer/assets/chall/view.html",
    }
    scripts = {
        "create": "/plugins/ctfd-instancer/assets/chall/create.js",
        "update": "/plugins/ctfd-instancer/assets/chall/update.js",
        "view": "/plugins/ctfd-instancer/assets/chall/view.js",
    }
    route = "/plugins/ctfd-instancer/assets/"
    blueprint = Blueprint(
        "instanced", __name__, template_folder="templates", static_folder="assets"
    )
    challenge_model = CTFdInstancedChallenge

    @staticmethod
    def instancetype(t):
        if t == "web": return 0
        if t == "socket": return 1
        raise ChallengeCreateException(f"Invalid challenge type: {t}")

    @classmethod
    def create(cls, request):
        data = request.form or request.get_json()

        challenge = cls.challenge_model(**data)

        db.session.add(challenge)
        db.session.commit()

        try:
            api_base = get_instancer_api_url()
            api_url = f"{api_base}/api/v1/challenge/"
            api_auth = {"Authorization": get_instancer_api_token()}
            payload = {
                "name": challenge.name,
                "description": challenge.description,
                "category": challenge.category,
                "type": CTFdInstancedTypeChallenge.instancetype(challenge.challenge_type),
                "remote_id": str(challenge.id),
                "flag": challenge.flag_base,
                "flag_type": challenge.flag_type,
                "duration": challenge.duration,
                "repository": challenge.repository,
                "chart": challenge.chart,
                "chart_version": challenge.chart_version,
                "values": challenge.values,
            }
            r = requests.post(api_url, json=payload, headers=api_auth, timeout=5)
            if r.status_code != 201:
                raise ChallengeCreateException(f"Failed to create challenge on Instancer: {r.text}")

            challenge.instancer_id = r.json().get("id")
            db.session.commit()
        except Exception as e:
            db.session.delete(challenge)
            db.session.commit()
            raise ChallengeCreateException(f"Error creating challenge: {str(e)}")

        return challenge

    @classmethod
    def read(cls, challenge):
        data = super().read(challenge)
        return data

    @classmethod
    def update(cls, challenge, request):
        data = request.form or request.get_json()

        for attr, value in data.items():
            if attr in ("flag_type", "duration"):
                try:
                    value = int(value)
                except (ValueError, TypeError):
                    raise ChallengeUpdateException(f"Invalid input for '{attr}'")
            setattr(challenge, attr, value)

        try:
            api_base = get_instancer_api_url()
            api_url = f"{api_base}/api/v1/challenge/{challenge.instancer_id}"
            api_auth = {"Authorization": get_instancer_api_token()}
            payload = {
                "name": challenge.name,
                "description": challenge.description,
                "category": challenge.category,
                "type": CTFdInstancedTypeChallenge.instancetype(challenge.challenge_type),
                "remote_id": str(challenge.id),
                "flag": challenge.flag_base,
                "flag_type": challenge.flag_type,
                "duration": challenge.duration,
                "repository": challenge.repository,
                "chart": challenge.chart,
                "chart_version": challenge.chart_version,
                "values": challenge.values,
            }
            r = requests.put(api_url, json=payload, headers=api_auth, timeout=5)
            if r.status_code != 200:
                raise ChallengeUpdateException(f"Failed to fetch challenge from Instancer: {r.text}")
        except Exception as e:
            raise ChallengeUpdateException(f"Error reading challenge data: {str(e)}")

        db.session.commit()
        return challenge

    @classmethod
    def delete(cls, challenge):
        if challenge.instancer_id is not None:
            try:
                api_base = get_instancer_api_url()
                api_url = f"{api_base}/api/v1/challenge/{challenge.instancer_id}"
                api_auth = {"Authorization": get_instancer_api_token()}
                r = requests.delete(api_url, headers=api_auth, timeout=5)
                if r.status_code != 204:
                    raise ChallengeUpdateException(f"Failed to delete challenge on Instancer: {r.text}")

            except Exception as e:
                raise ChallengeUpdateException(f"Error deleting challenge: {str(e)}")

        super().delete(challenge)

    @classmethod
    def solve(cls, user, team, challenge, request):
        super().solve(user, team, challenge, request)
