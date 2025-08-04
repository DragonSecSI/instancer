from flask import Blueprint, render_template, request
from CTFd.utils.decorators import admins_only
from CTFd.plugins import bypass_csrf_protection

from .config import get_instancer_api_url, get_instancer_api_token, set_instancer_api_url, set_instancer_api_token, get_instancer_public_url, set_instancer_public_url

instancer_admin_bp = Blueprint(
    "instancer_admin", __name__, template_folder="assets"
)

@instancer_admin_bp.route("/admin/plugins/instancer", methods=["GET", "POST"])
@admins_only
@bypass_csrf_protection
def instancer_settings():
    saved = False
    if request.method == "POST":
        url = request.form.get("INSTANCER_API_URL", "").strip()
        token = request.form.get("INSTANCER_API_TOKEN", "").strip()
        public = request.form.get("INSTANCER_PUBLIC_URL", "").strip()
        set_instancer_api_url(url)
        set_instancer_api_token(token)
        set_instancer_public_url(public)
        saved = True

    return render_template(
        "admin/settings.html",
        instancer_api_url=get_instancer_api_url(),
        instancer_api_token=get_instancer_api_token(),
        instancer_public_url=get_instancer_public_url(),
        saved=saved,
    )

