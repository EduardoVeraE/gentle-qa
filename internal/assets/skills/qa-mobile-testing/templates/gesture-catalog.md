<!-- Skill: qa-mobile-testing · Template: gesture-catalog -->
<!-- Placeholders: {{element_id}}, {{target_element_id}} -->

# Gesture Catalog

Cross-framework reference for the most common mobile gestures. Each entry lists when to use the gesture and minimal code for **Appium (TypeScript / WebdriverIO)**, **Detox (JavaScript)**, **XCUITest (Swift)**, and **Espresso (Kotlin)**.

> `{{element_id}}` is the accessibility id (`testID` in RN, `accessibilityIdentifier` on iOS, `contentDescription` / `resource-id` on Android).

## 1. Tap
When: primary activation — buttons, list items, links.
```ts
await driver.$(`~{{element_id}}`).click(); // Appium
```
```js
await element(by.id('{{element_id}}')).tap(); // Detox
```
```swift
app.buttons["{{element_id}}"].tap() // XCUITest
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(click()) // Espresso
```

## 2. Double Tap
When: zoom toggle on maps/images, like-on-double-tap.
```ts
await driver.$(`~{{element_id}}`).doubleClick();
```
```js
await element(by.id('{{element_id}}')).multiTap(2);
```
```swift
app.images["{{element_id}}"].doubleTap()
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(doubleClick())
```

## 3. Long Press
When: context menus, drag initiation, reorder handle.
```ts
await driver.$(`~{{element_id}}`).touchAction([{ action: 'longPress' }, { action: 'release' }]);
```
```js
await element(by.id('{{element_id}}')).longPress(1000);
```
```swift
app.cells["{{element_id}}"].press(forDuration: 1.0)
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(longClick())
```

## 4. Swipe (4 directions)
When: navigation between pages, dismiss cards, reveal row actions. Directions: `up | down | left | right`.
```ts
await (await driver.$(`~{{element_id}}`)).swipe('left');
```
```js
await element(by.id('{{element_id}}')).swipe('left', 'fast', 0.75);
```
```swift
app.otherElements["{{element_id}}"].swipeLeft() // .swipeRight/.swipeUp/.swipeDown
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(swipeLeft()) // swipeRight/swipeUp/swipeDown
```

## 5. Pinch In / Out
When: zoom on images, maps, PDFs. `scale > 1` zooms in, `< 1` zooms out.
```ts
await driver.execute('mobile: pinch', { element: (await driver.$(`~{{element_id}}`)).elementId, scale: 0.5, velocity: -1.0 });
```
```js
await element(by.id('{{element_id}}')).pinch(0.75);
```
```swift
app.maps["{{element_id}}"].pinch(withScale: 2.0, velocity: 1.0)
```
```kotlin
// Espresso lacks pinch — use UiAutomator2:
UiDevice.getInstance(InstrumentationRegistry.getInstrumentation())
  .findObject(By.res("{{element_id}}")).pinchClose(0.75f)
```

## 6. Rotate (two-finger rotate)
When: rotating an image canvas or map heading.
```ts
await driver.execute('mobile: rotateGesture', { elementId: (await driver.$(`~{{element_id}}`)).elementId, rotation: 1.57 });
```
```js
// Detox has no native rotate — expose via custom action
await element(by.id('{{element_id}}')).performAccessibilityAction('rotate');
```
```swift
app.otherElements["{{element_id}}"].rotate(.pi / 2, withVelocity: 1.0)
```
```kotlin
// UiAutomator2 substitute
device.findObject(By.res("{{element_id}}")).pinchOpen(0.5f)
```

## 7. Drag and Drop
When: reorder lists, move cards across columns.
```ts
const src = await driver.$(`~{{element_id}}`);
const dst = await driver.$(`~{{target_element_id}}`);
await driver.execute('mobile: dragGesture', { elementId: src.elementId, endX: (await dst.getLocation()).x, endY: (await dst.getLocation()).y });
```
```js
await element(by.id('{{element_id}}')).longPressAndDrag(2000, NaN, NaN, element(by.id('{{target_element_id}}')));
```
```swift
app.cells["{{element_id}}"].press(forDuration: 1.0, thenDragTo: app.cells["{{target_element_id}}"])
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(
  GeneralSwipeAction(Swipe.SLOW, GeneralLocation.CENTER, GeneralLocation.BOTTOM_CENTER, Press.FINGER))
```

## 8. Scroll To Element
When: list is virtualized and the target is offscreen.
```ts
await driver.execute('mobile: scroll', { strategy: 'accessibility id', selector: '{{element_id}}' });
```
```js
await waitFor(element(by.id('{{element_id}}'))).toBeVisible()
  .whileElement(by.id('list-{{element_id}}')).scroll(200, 'down');
```
```swift
let cell = app.cells["{{element_id}}"]; while !cell.isHittable { app.swipeUp() }; cell.tap()
```
```kotlin
onView(withId(R.id.list)).perform(
  RecyclerViewActions.scrollTo<RecyclerView.ViewHolder>(hasDescendant(withText("{{element_id}}"))))
```

## 9. Pull to Refresh
When: list refresh interaction.
```ts
await (await driver.$(`~{{element_id}}`)).touchAction([
  { action: 'press', x: 200, y: 200 }, { action: 'moveTo', x: 200, y: 600 }, { action: 'release' }]);
```
```js
await element(by.id('{{element_id}}')).swipe('down', 'slow', 0.9);
```
```swift
let list = app.collectionViews["{{element_id}}"]
list.coordinate(withNormalizedOffset: CGVector(dx: 0.5, dy: 0.1))
  .press(forDuration: 0.05, thenDragTo: list.coordinate(withNormalizedOffset: CGVector(dx: 0.5, dy: 0.9)))
```
```kotlin
onView(withId(R.id.{{element_id}})).perform(swipeDown())
```

## 10. Two-Finger Tap
When: secondary actions (e.g. zoom-out on maps, accessibility shortcuts).
```ts
await driver.execute('mobile: twoFingerTap', { elementId: (await driver.$(`~{{element_id}}`)).elementId });
```
```js
await element(by.id('{{element_id}}')).multiTap(2); // approximate; Detox lacks true two-finger
```
```swift
app.maps["{{element_id}}"].twoFingerTap()
```
```kotlin
device.findObject(By.res("{{element_id}}")).click() // UiAutomator2 fallback
```

---

## Notes
- iOS Simulator cannot synthesize Force / Haptic Touch — use real devices.
- Espresso is in-process and cannot synthesize multi-touch — use UiAutomator2 for pinch / rotate / multi-finger.
- Detox is synchronous-by-design — gestures queue; do not insert manual sleeps.
- Always prefer accessibility ids over coordinates — stable across devices, screen sizes, and locales.
