import * as userService from '../../api/user/user.service';
import * as groupService from '../../api/group/group.service';
import { create } from "zustand"
import useAuthStore from '../auth';

interface IGlobalStore {
    /* Conversations > people */
    people: any[]
    fetchPeople: () => void

    /* Conversations > groups */
    groups: any[]
    fetchGroups: () => void

}

const useGlobalStore = create<IGlobalStore>((set, get) => ({
    people: [],
    fetchPeople: async () => {
        try {
            const response = await userService.getUsers()
            const { user: me } = useAuthStore.getState()

            set({ people: response.data.filter((user: any) => user.id !== me?.id) })
        } catch (error: any) {
            console.log(error)

            set({ people: [] })
        }
    },
    
    groups: [],
    fetchGroups: async () => {
        try {
            const response = await groupService.getGroups()

            set({ groups: response.data })
        } catch (error: any) {
            console.log(error)

            set({ groups: [] })
        }
    },
}));

export default useGlobalStore