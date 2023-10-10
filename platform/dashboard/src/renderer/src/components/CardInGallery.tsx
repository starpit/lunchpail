import { PureComponent } from "react"
import { Card, CardHeader, CardTitle, CardBody, CardFooter } from "@patternfly/react-core"

import type { MouseEvent, ReactNode } from "react"
import type { CardHeaderActionsObject } from "@patternfly/react-core"

import { dl } from "./DescriptionGroup"

import type { Kind } from "../names"
import type { DrilldownProps } from "../context/DrawerContext"

import "./CardInGallery.scss"

type BaseProps = DrilldownProps

export default abstract class CardInGallery<Props> extends PureComponent<Props & BaseProps> {
  protected readonly stopPropagation = (evt: MouseEvent<HTMLElement>) => evt.stopPropagation()

  protected abstract kind(): Kind

  protected abstract label(): string

  protected abstract icon(): ReactNode

  /** DescriptionList groups to display in the Card summary */
  protected abstract groups(): ReactNode[]

  protected actions(): undefined | CardHeaderActionsObject {
    return undefined
  }

  private readonly onClick = () => {
    this.props.showDetails({ id: this.label(), kind: this.kind() })
  }

  private header() {
    return (
      <CardHeader actions={this.actions()} className="codeflare--card-header-no-wrap">
        <span className="codeflare--card-icon">{this.icon()}</span>
      </CardHeader>
    )
  }

  private title() {
    return this.label()
  }

  private body() {
    return dl(this.groups(), { isCompact: true })
  }

  protected footer(): null | ReactNode {
    return null
  }

  private card() {
    return (
      <Card
        isClickable
        isSelectable
        isSelectableRaised
        isSelected={this.props.currentlySelectedId === this.label() && this.props.currentlySelectedKind === this.kind()}
        onClick={this.onClick}
      >
        {this.header()}
        <CardTitle>{this.title()}</CardTitle>
        <CardBody>{this.body()}</CardBody>
        {this.footer() && <CardFooter>{this.footer()}</CardFooter>}
      </Card>
    )
  }

  public override render() {
    return this.card()
  }
}