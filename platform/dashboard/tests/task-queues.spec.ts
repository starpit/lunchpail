// @ts-check
import { expect, test } from "@playwright/test"
import launchElectron from "./launch-electron"

test("3 demo task queues visible", async () => {
  // Launch Electron app.
  const electronApp = await launchElectron()

  // Get the first page that the app opens, wait if necessary.
  const page = await electronApp.firstWindow()

  // Check if we are in demo mode (should be true by default)
  const demoModeStatus = await page.getByLabel("Demo Mode").isChecked()
  console.log(`Demo mode on?: ${demoModeStatus}`)

  // If in demo mode, then continue with Task queue card test
  if (demoModeStatus) {
    // Get Task Queue tab element from the sidebar and click
    await page.getByRole("link", { name: "Task Queues" }).click()

    // Verify that the three showing are the pink, purple, and green cards
    const expectedCards = ["green", "pink", "purple"]

    await Promise.all(
      expectedCards.map((id) =>
        expect(page.locator(`[data-ouia-component-type="PF5/Card"][data-ouia-component-id="${id}"]`))
          .toBeVisible({
            timeout: 30000,
          })
          .then(() => console.log("got", id)),
      ),
    )
  }
})