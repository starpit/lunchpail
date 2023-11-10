import None from "@jay/components/None"
import CardInGallery from "@jay/components/CardInGallery"
import { linkToAllDetails } from "@jay/renderer/navigate/details"
import { descriptionGroup } from "@jay/components/DescriptionGroup"

import { name as datasetsName } from "../../datasets/name"
// import { name as taskqueuesName } from "../../taskqueues/name"

import type Props from "./Props"

import ApplicationIcon from "./Icon"

export function api(props: Props) {
  const { api } = props.application.spec

  if (api === "workqueue") {
    return []
  } else {
    return [descriptionGroup("api", api, undefined, "The API used by this Application to distribute work.")]
  }
}

function inputs(props: Props) {
  return props.application.spec.inputs
    ? props.application.spec.inputs.flatMap((_) => Object.values(_.sizes)).filter(Boolean)
    : []
}

export function taskqueues(props: Props) {
  return inputs(props).filter(
    (taskqueueName) => !!props.taskqueues.find((taskqueue) => taskqueueName === taskqueue.metadata.name),
  )
}

function datasets(props: Props) {
  return inputs(props).filter(
    (datasetName) => !!props.datasets.find((dataset) => datasetName === dataset.metadata.name),
  )
}

/* export function taskqueuesGroup(props: Props) {
  const queues = taskqueues(props)

  return (
    queues.length > 0 &&
    descriptionGroup(
      taskqueuesName,
      queues.length === 0 ? None() : linkToAllDetails("taskqueues", queues),
      queues.length,
      "The Task Queues this application is capable of processing, i.e. those that it is compatible with.",
    )
  )
} */

export function datasetsGroup(props: Props) {
  const data = datasets(props)

  return (
    data.length > 0 &&
    descriptionGroup(
      datasetsName,
      data.length === 0 ? None() : linkToAllDetails("datasets", data),
      data.length,
      "The Datasets this application requires as input.",
    )
  )
}

function hasWorkerPool(props: Props) {
  return !!props.workerpools.find((_) => _.spec.application.name === props.application.metadata.name)
}

export default function ApplicationCard(props: Props) {
  const icon = <ApplicationIcon application={props.application} hasWorkerPool={hasWorkerPool(props)} />
  const name = props.application.metadata.name

  const groups = [
    ...api(props),
    props.application.spec.description && descriptionGroup("Description", props.application.spec.description),
    // taskqueuesGroup(props),
    datasetsGroup(props),
  ]

  return <CardInGallery kind="applications" name={name} icon={icon} groups={groups} />
}