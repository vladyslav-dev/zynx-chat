import useAuthStore from "../../../shared/store/auth";
import { Box, Button, Flex, Heading, Link, Text, TextField } from "@radix-ui/themes";

const Login = () => {
    const { login } = useAuthStore()

    const onSubmit = async (e: React.FormEvent) => {
        e.preventDefault()

        const formData = new FormData(e.target as HTMLFormElement);

        const phone = formData.get('phone') as string
        const password = formData.get('password') as string
        
        login({ phone, password })
    }

    return (
        <Flex width={"100%"} height={"100%"} direction={"column"} justify={"center"} align={"center"}>
            <Heading as={"h2"} className="mb-4">Login</Heading>
            <form onSubmit={onSubmit}>
                <Flex direction="column" gap={"16px"}>
                    <Box maxWidth="220px">
                        <Text as="label" htmlFor="phone">Phone</Text>
                        <TextField.Root size="2" placeholder="Phone" id="phone" name="phone" />
                    </Box>
                    <Box maxWidth="220px">
                        <Text as="label" htmlFor="password">Password</Text>
                        <TextField.Root size="2" placeholder="Password" id="password" name="password" type="password" />
                    </Box>
                    <Flex maxWidth="220px" gap={"8px"} direction={"column"} justify={"center"}>
                        <Button type="submit" className="w-full">Login</Button>
                        <Link href={"/register"} size={"1"}>Create account?</Link>
                    </Flex>
                </Flex>
            </form>
            </Flex>
    )
}

export default Login