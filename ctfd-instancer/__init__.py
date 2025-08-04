from CTFd.plugins import register_plugin_assets_directory, register_user_page_menu_bar
from CTFd.plugins.challenges import CHALLENGE_CLASSES
from CTFd.plugins.flags import FLAG_CLASSES
from CTFd.models import db

from .flag import InstancedFlag, InstancedFlagAudit
from .challenge import instanced_challenge_bp, CTFdInstancedChallenge, CTFdInstancedTypeChallenge
from .hooks import instancer_bp, InstancerTokenTable
from .admin import instancer_admin_bp

def load(app):
    with app.app_context():
        db.create_all()

    register_plugin_assets_directory(app, base_path='/plugins/ctfd-instancer/assets/')
    
    FLAG_CLASSES["instanced"] = InstancedFlag
    CHALLENGE_CLASSES["instanced"] = CTFdInstancedTypeChallenge
    
    app.register_blueprint(instancer_bp)
    app.register_blueprint(instanced_challenge_bp)
    app.register_blueprint(instancer_admin_bp)
    
    register_user_page_menu_bar("Instancer", "/instancer")
    
    app.logger.info("CTFd Instancer plugin loaded.")
