agent_list = dict()


def update_agents(new_agent_list):
    # delete old one
    del_list = []
    for agent in agent_list:
        if len(list(filter(lambda new_agent: new_agent["name"] == agent_list[agent]["name"], new_agent_list))) == 0:
            del_list.append(agent_list[agent]["name"])
    for del_l in del_list:
        del (agent_list[del_l])

    # update new one
    for new_agent in new_agent_list:  # dict list
        if not new_agent["name"] in agent_list:
            agent_list[new_agent["name"]] = new_agent


def get_all_agents():
    return agent_list


def get_agent(name):
    return agent_list[name]
