import { ReactElement } from "react"
import gear from "../../../assets/icons/gear.svg"

const GearIcon = (props: React.ImgHTMLAttributes<HTMLImageElement>): ReactElement => <img src={gear} alt='gear' {...props} />

export default GearIcon