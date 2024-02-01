import type { ContentProviderSidebarSpec } from "./ContentProvider"

export const Component = "Component"
export const Components = "Components"

/** Sidebar Group model for Job-related resources */
export function componentsSidebar(priority: number) {
  return {
    priority,
    group: Components,
  } satisfies ContentProviderSidebarSpec
}

/** Sidebar Group model for Secrets-related resources */
export const secretsSidebar = {
  group: "Secrets",
} satisfies ContentProviderSidebarSpec