import { Box, Flex } from "@radix-ui/themes"
import { ReactNode } from "react"
import Navigation from "../components/Navigation"

const MainLayout = ({ children }: { children: ReactNode }) => {
    return (
        <Flex>
            <Navigation />
            {children}
        </Flex>
    )
}

export default MainLayout