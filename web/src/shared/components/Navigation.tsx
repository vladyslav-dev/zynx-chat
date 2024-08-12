import { Box, Flex, Popover, Link, Text } from "@radix-ui/themes"
import ChatIcon from "./icons/ChatIcon"
import GearIcon from "./icons/Gear"
import useAuthStore from "../store/auth"
import { NavLink } from "react-router-dom"

const Navigation = () => {
    const { logout, user: me } = useAuthStore()

    return (
        <Flex direction={"column"} justify={"between"} height={"100dvh"} className="bg-slate-100">
            <Flex direction={"column"}>
                <NavLink to={"/"} className={({ isActive, isPending }) => [
                    isActive ? "bg-indigo-200" : "",
                    "p-4"
                ].join(" ")}>
                    <Flex direction={"column"} gap={"16px"} align={"center"}>
                        <ChatIcon className="w-[24px]" />
                        <Text size={"1"} as="label">Chats</Text>
                    </Flex>
                </NavLink>
                <NavLink to={"/settings"} className={({ isActive, isPending }) => [
                    isActive ? "bg-indigo-200" : "",
                    "p-4"
                ].join(" ")}>
                    <Flex direction={"column"} gap={"16px"} align={"center"}>
                        <GearIcon className="w-[24px]" />
                        <Text size={"1"} as="label">Settings</Text>
                    </Flex>
                </NavLink>
            </Flex>
            <Flex align={"center"} direction={"column"}>
                <Popover.Root>
                    <Popover.Trigger className="cursor-pointer">
                        <Box className="p-4">
                            <Text size={"2"}>{me?.username}</Text>
                        </Box>
                    </Popover.Trigger>
                    <Popover.Content side={"right"} align={"start"}>
                        <Box className="p-4">
                            <NavLink to={"/profile"}>
                                <Link size={"2"}>Profile</Link>
                            </NavLink>
                        </Box>
                        <Box className="p-4">
                            <Link size={"2"} onClick={logout} className="text-[red] cursor-pointer">Logout</Link>
                        </Box>
                    </Popover.Content>
                </Popover.Root>
            </Flex>
        </Flex>
    )
}

export default Navigation