import $api from "../index"

export const getUsers = async () => {
    return $api.get("/user/getAll")
}

export const getUsersByIds = async (ids: number[]) => {
    return $api.post("/user/getByIds", ids)
}

export const getUsersByGroupId = async (group_id: number) => {
    return $api.get(`/user/getByGroupId?group_id=${group_id}`)
}