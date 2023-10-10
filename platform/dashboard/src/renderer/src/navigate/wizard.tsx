import { Link } from "react-router-dom"
import { Button } from "@patternfly/react-core"

import { stopPropagation } from "."

import RocketIcon from "@patternfly/react-icons/dist/esm/icons/rocket-icon"
import PlusCircleIcon from "@patternfly/react-icons/dist/esm/icons/plus-circle-icon"

type StartOrAdd = "start" | "add" | "create"

function href(view: string, returnTo?: string, hash?: string, qs: string[] = []) {
  const queries = [`view=${view}`, ...qs, returnTo ? `returnTo=${returnTo}` : undefined].filter(Boolean)

  return "?" + queries.join("&") + (hash ?? "")
}

/** A React component that will offer a Link to a given `data-href` */
function linker(props: { "data-href": string; "data-link-text": string; "data-start-or-add": StartOrAdd }) {
  const href = props["data-href"]
  const start = props["data-start-or-add"]

  const icon = start === "start" ? <RocketIcon /> : <PlusCircleIcon />
  const linkText = props["data-link-text"]

  return (
    <Link {...props} to={href}>
      <span className="pf-v5-c-button__icon pf-m-start">{icon}</span> {linkText}
    </Link>
  )
}

/** Base/public props for subclasses */
export type WizardProps = Omit<import("../router/withLocation").LocationProps, "navigate"> & {
  startOrAdd?: StartOrAdd
}

/** Internal props */
type Props = WizardProps & {
  view: string
  linkText: string
  qs: string[]
}

/**
 * @return a UI component that links to the a wizard `view`. If
 * `startOrAdd` is `start`, then present the UI as if this were the
 * first time we were asking to create such a thing via a wizard;
 * otherwise, present as if we are augmenting an existing thing.
 */
export default function LinkToNewWizard(props: Props) {
  const currentHash = props.location.hash
  const currentSearch = props.searchParams
  const returnTo = encodeURIComponent(`?${currentSearch}`)
  const theHref = href(props.view, returnTo, currentHash, props.qs)

  return (
    <Button
      isInline={props.startOrAdd === "create"}
      variant={props.startOrAdd === "create" ? "link" : "primary"}
      size="sm"
      onClick={stopPropagation}
      data-start-or-add={props.startOrAdd || "start"}
      data-link-text={props.linkText}
      data-href={theHref}
      component={linker}
    />
  )
}