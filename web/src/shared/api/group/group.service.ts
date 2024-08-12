import $api from "../index"

export const getGroupById = async (id: number) => {
    return $api.get(`/group/get?id=${id}`)
}

export const getGroups = async () => {
    return $api.get("/group/getAll")
}