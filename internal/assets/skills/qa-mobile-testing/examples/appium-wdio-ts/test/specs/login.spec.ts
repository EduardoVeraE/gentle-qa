import { expect } from '@wdio/globals';
import LoginPage from '../pages/login.page.js';

describe('Login flow', () => {
  beforeEach(async () => {
    await LoginPage.open();
  });

  it('logs in with valid credentials', async () => {
    await LoginPage.login('demo@example.com', 'correct-horse-battery-staple');
    await expect(LoginPage.homeHeader).toBeDisplayed();
  });

  it('shows an error for invalid credentials', async () => {
    await LoginPage.login('demo@example.com', 'wrong');
    await expect(LoginPage.errorMessage).toBeDisplayed();
    await expect(LoginPage.errorMessage).toHaveText(
      expect.stringContaining('Invalid')
    );
  });
});
