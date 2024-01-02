from tinydb import TinyDB, Query
db = TinyDB('db.json')
db.drop_tables()
HPA = Query()


def update_hpas(hpas):
    # delete old one
    data = db.all()
    for d in data:
        if len(list(filter(lambda hpa: hpa["name"] == d["name"], hpas))) == 0:
            db.remove(HPA.name == d["name"])
    # update new one
    for hpa in hpas:
        if not db.contains(HPA.name == hpa["name"]):
            db.insert(hpa)


def get_all_hpas():
    return db.all()


def get_hpa(name):
    return db.search(HPA.name == name)
