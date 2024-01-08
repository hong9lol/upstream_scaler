hpa_list = dict()


def update_hpas(new_hpa_list):
    # delete old one
    del_list = []
    for hpa in hpa_list:
        if len(list(filter(lambda new_hpa: new_hpa["name"] == hpa_list[hpa]["name"], new_hpa_list))) == 0:
            del_list.append(hpa_list[hpa]["name"])
    for del_l in del_list:
        del (hpa_list[del_l])

    # update new one
    for new_hpa in new_hpa_list:  # dict list
        if not new_hpa["name"] in hpa_list:
            hpa_list[new_hpa["name"]] = new_hpa


def get_all_hpas():
    return hpa_list


def get_hpa(name):
    return hpa_list[name]


# from tinydb import TinyDB, Query
# db = TinyDB('db.json')
# db.drop_tables()
# HPA = Query()


# def update_hpas(hpas):
#     # delete old one
#     data = db.all()
#     for d in data:
#         if len(list(filter(lambda hpa: hpa["name"] == d["name"], hpas))) == 0:
#             db.remove(HPA.name == d["name"])
#     # update new one
#     for hpa in hpas:
#         if not db.contains(HPA.name == hpa["name"]):
#             db.insert(hpa)


# def get_all_hpas():
#     return db.all()


# def get_hpa(name):
#     return db.search(HPA.name == name)
