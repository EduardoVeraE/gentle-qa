describe('App smoke', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
  });

  beforeEach(async () => {
    await device.reloadReactNative();
  });

  it('launches and shows the welcome screen', async () => {
    await expect(element(by.id('welcome-screen'))).toBeVisible();
  });

  it('navigates to the login screen', async () => {
    await element(by.id('go-to-login')).tap();
    await expect(element(by.id('login-screen'))).toBeVisible();
  });
});
