hpa_list = dict()


def update_hpas(new_hpa_list):
    # delete old one
    del_list = []
    for hpa in hpa_list:
        if (
            len(
                list(
                    filter(
                        lambda new_hpa: new_hpa["name"] == hpa_list[hpa]["name"],
                        new_hpa_list,
                    )
                )
            )
            == 0
        ):
            del_list.append(hpa_list[hpa]["name"])
    for del_l in del_list:
        del hpa_list[del_l]

    # update new one
    for new_hpa in new_hpa_list:  # dict list
        if not new_hpa["name"] in hpa_list:
            hpa_list[new_hpa["name"]] = new_hpa

    # maybe just... TODO need to check it works well
    # hpa_list.clear()
    # hpa_list.update(new_hpa_list)


def get_all_hpas():
    return hpa_list


def get_hpa(name):
    return hpa_list[name]
