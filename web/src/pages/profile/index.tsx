import useAuthStore from "../../shared/store/auth"
import { Box, Flex, Heading, Text } from "@radix-ui/themes"

const Profile = () => {
    const { user: me } = useAuthStore()

    if (!me) {
        return null
    }

    return (
        <Box className="p-4">
            <Heading as="h1" style={{ marginBottom: "1rem" }}>
                Profile
            </Heading>

            <Flex direction="column" gap="0.5rem">
                <Flex direction="row">
                    <Text as="label" style={{ fontWeight: "bold" }}>
                        ID:
                    </Text>
                    <Text>{me.id}</Text>
                </Flex>
                <Flex direction="row">
                    <Text as="label" style={{ fontWeight: "bold" }}>
                        Username:
                    </Text>
                    <Text>{me.username}</Text>
                </Flex>
                <Flex direction="row">
                    <Text as="label" style={{ fontWeight: "bold" }}>
                        Phone:
                    </Text>
                    <Text>{me.phone}</Text>
                </Flex>
            </Flex>
        </Box>
    )
}

export default Profile