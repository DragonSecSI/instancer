from flask import current_app
from CTFd.models import Configs, db

def _get_config_value(key, default=None):
    cfg = Configs.query.filter_by(key=key).first()
    if cfg and cfg.value:
        return cfg.value
    return current_app.config.get(key, default)


def _set_config_value(key, value):
    cfg = Configs.query.filter_by(key=key).first()
    if cfg:
        cfg.value = value
    else:
        cfg = Configs(key=key, value=value)
        db.session.add(cfg)
    db.session.commit()


def get_instancer_api_url():
    return _get_config_value("INSTANCER_API_URL", None)

def set_instancer_api_url(value):
    _set_config_value("INSTANCER_API_URL", value)


def get_instancer_api_token():
    return _get_config_value("INSTANCER_API_TOKEN", None)

def set_instancer_api_token(value):
    _set_config_value("INSTANCER_API_TOKEN", value)


def get_instancer_public_url():
    return _get_config_value("INSTANCER_PUBLIC_URL", None)

def set_instancer_public_url(value):
    _set_config_value("INSTANCER_PUBLIC_URL", value)
