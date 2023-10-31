import type TaskQueueEvent from "@jay/common/events/TaskQueueEvent"
import type EventSourceLike from "@jay/common/events/EventSourceLike"

import Base from "./base"
import { ns } from "./misc"

export const colors = ["pink", "green", "purple"]

export function inbox(taskqueue: TaskQueueEvent) {
  return parseInt(taskqueue.metadata.annotations["codeflare.dev/unassigned"] || "0", 10)
}

export function inboxIncr(taskqueue: TaskQueueEvent, incr = 1) {
  taskqueue.metadata.annotations["codeflare.dev/unassigned"] = String(inbox(taskqueue) + incr)
}

export default class DemoTaskQueueEventSource extends Base implements EventSourceLike {
  private readonly endpoints = ["e1", "e2", "e3"]
  private readonly buckets = ["pile1", "pile2", "pile3"]
  private readonly isReadOnly = [true, false, true]

  private readonly taskqueues: TaskQueueEvent[] = Array(3)
    .fill(0)
    .map((_, idx) => ({
      metadata: {
        name: colors[idx],
        namespace: ns,
        creationTimestamp: new Date().toUTCString(),
        annotations: {
          "codeflare.dev/status": "Ready",
          "codeflare.dev/unassigned": "0",
        },
      },
      spec: {
        idx,
        local: {
          type: "COS",
          endpoint: this.endpoints[idx],
          bucket: this.buckets[idx],
          readonly: this.isReadOnly[idx],
        },
      },
    }))

  public get sets(): readonly Omit<TaskQueueEvent, "timestamp">[] {
    return this.taskqueues
  }

  private sendEventFor = (
    taskqueue: (typeof this.taskqueues)[number],
    status = taskqueue.metadata.annotations["codeflare.dev/status"],
  ) => {
    const model: TaskQueueEvent = Object.assign({}, taskqueue, {
      status,
      timestamp: Date.now(),
      //inbox: ~~(Math.random() * 20),
      //outbox: ~~(Math.random() * 2),
    })
    this.handlers.forEach((handler) => handler(new MessageEvent("taskqueue", { data: JSON.stringify([model]) })))
  }

  protected override initInterval(intervalMillis: number) {
    if (!this.interval) {
      const { taskqueues, sendEventFor } = this

      this.interval = setInterval(
        (function interval() {
          const whichToUpdate = Math.floor(Math.random() * taskqueues.length)
          const taskqueue = taskqueues[whichToUpdate]
          inboxIncr(taskqueue)
          sendEventFor(taskqueue)
          return interval
        })(), // () means invoke the interval right away
        intervalMillis,
      )
    }
  }

  public override delete(props: { name: string; namespace: string }) {
    const idx = this.taskqueues.findIndex(
      (_) => _.metadata.name === props.name && _.metadata.namespace === props.namespace,
    )
    if (idx >= 0) {
      const model = this.taskqueues[idx]
      this.taskqueues.splice(idx, 1)
      this.sendEventFor(model, "Terminating")
      return true
    } else {
      return {
        code: 404,
        message: "Resource not found",
      }
    }
  }
}