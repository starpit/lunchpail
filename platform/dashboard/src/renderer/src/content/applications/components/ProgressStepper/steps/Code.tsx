import { Link } from "react-router-dom"

import type Step from "../Step"
import { repoPlusSource } from "../../tabs/Code"

const step: Step = {
  id: "Code",
  variant: () => "success",
  content: (props, onClick) => (
    <span>
      Code will be pulled from{" "}
      <Link onClick={onClick} target="_blank" to={props.application.spec.repo}>
        {repoPlusSource(props)}
      </Link>
    </span>
  ),
}

export default step