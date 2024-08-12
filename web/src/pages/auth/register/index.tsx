import { Box, Button, Flex, Heading, Link, Text, TextField } from "@radix-ui/themes";
import useAuthStore from "../../../shared/store/auth";


const Register = () => {
    const { register } = useAuthStore()

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const username = formData.get('username') as string
        const phone = formData.get('phone') as string
        const password = formData.get('password') as string
        
        register({ username, phone, password })
    }

    return (
        <Flex width={"100%"} height={"100%"} direction={"column"} justify={"center"} align={"center"}>
            <Heading as={"h2"} className="mb-4">Register</Heading>
            <form onSubmit={onSubmit}>
                <Flex direction="column" gap={"16px"}>
                    <Box maxWidth="220px">
                        <Text as="label" htmlFor="username">Username</Text>
                        <TextField.Root size="2" placeholder="Username" id="username" name="username" />
                    </Box>
                    <Box maxWidth="220px">
                        <Text as="label" htmlFor="phone">Phone</Text>
                        <TextField.Root size="2" placeholder="Phone" id="phone" name="phone" />
                    </Box>
                    <Box maxWidth="220px">
                        <Text as="label" htmlFor="password">Password</Text>
                        <TextField.Root size="2" placeholder="Password" id="password" name="password" type="password" />
                    </Box>
                    <Flex maxWidth="220px" gap={"8px"} direction={"column"} justify={"center"}>
                        <Button type="submit" className="w-full">Register</Button>
                        <Link href={"/login"} size={"1"}>Already have an account?</Link>
                    </Flex>
                </Flex>
            </form>
        </Flex>
    )
}

export default Register