import { browser, $ } from '@wdio/globals';

/**
 * PageObject for the Login screen.
 * Selectors use accessibility ids which work cross-platform via Appium.
 */
class LoginPage {
  public get usernameField() {
    return $('~login-username');
  }

  public get passwordField() {
    return $('~login-password');
  }

  public get submitButton() {
    return $('~login-submit');
  }

  public get errorMessage() {
    return $('~login-error');
  }

  public get homeHeader() {
    return $('~home-header');
  }

  public async open(): Promise<void> {
    // The app launches at the login screen by default.
    // For deep-linking on iOS use `mobile: deepLink`.
    await browser.pause(500);
  }

  public async login(username: string, password: string): Promise<void> {
    await this.usernameField.waitForDisplayed();
    await this.usernameField.setValue(username);
    await this.passwordField.setValue(password);
    await this.submitButton.click();
  }
}

export default new LoginPage();
