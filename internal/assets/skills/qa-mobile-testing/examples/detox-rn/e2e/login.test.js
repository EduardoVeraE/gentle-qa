describe('Login flow', () => {
  beforeAll(async () => {
    await device.launchApp({ newInstance: true });
  });

  beforeEach(async () => {
    await device.reloadReactNative();
    await element(by.id('go-to-login')).tap();
    await expect(element(by.id('login-screen'))).toBeVisible();
  });

  it('logs in with valid credentials', async () => {
    await element(by.id('login-username')).typeText('demo@example.com');
    await element(by.id('login-password')).typeText('correct-horse-battery-staple');
    await element(by.id('login-submit')).tap();

    await waitFor(element(by.id('home-header')))
      .toBeVisible()
      .withTimeout(10000);
  });

  it('shows an error for invalid credentials', async () => {
    await element(by.id('login-username')).typeText('demo@example.com');
    await element(by.id('login-password')).typeText('wrong');
    await element(by.id('login-submit')).tap();

    await expect(element(by.id('login-error'))).toBeVisible();
    await expect(element(by.id('login-error'))).toHaveText(
      'Invalid credentials',
    );
  });
});
