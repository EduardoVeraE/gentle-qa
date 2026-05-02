import type { Options } from '@wdio/types';
import path from 'node:path';

const platform = (process.env.PLATFORM ?? 'ios').toLowerCase();

const iosCapabilities = [
  {
    platformName: 'iOS',
    'appium:automationName': 'XCUITest',
    'appium:platformVersion': process.env.IOS_VERSION ?? '17.4',
    'appium:deviceName': process.env.IOS_DEVICE ?? 'iPhone 15',
    'appium:app': path.resolve(process.cwd(), 'apps/MyApp.app'),
    'appium:newCommandTimeout': 240,
    'appium:wdaLaunchTimeout': 120000,
  },
];

const androidCapabilities = [
  {
    platformName: 'Android',
    'appium:automationName': 'UiAutomator2',
    'appium:platformVersion': process.env.ANDROID_VERSION ?? '14',
    'appium:deviceName': process.env.ANDROID_DEVICE ?? 'Pixel_8_API_34',
    'appium:app': path.resolve(process.cwd(), 'apps/app-debug.apk'),
    'appium:appWaitActivity': '*',
    'appium:autoGrantPermissions': true,
  },
];

export const config: Options.Testrunner = {
  runner: 'local',
  autoCompileOpts: {
    autoCompile: true,
    tsNodeOpts: {
      project: './tsconfig.json',
      transpileOnly: true,
    },
  },
  specs: ['./test/specs/**/*.spec.ts'],
  maxInstances: 1,
  capabilities: platform === 'android' ? androidCapabilities : iosCapabilities,
  logLevel: 'info',
  bail: 0,
  baseUrl: 'http://localhost',
  waitforTimeout: 15000,
  connectionRetryTimeout: 120000,
  connectionRetryCount: 3,
  services: [
    [
      'appium',
      {
        args: { basePath: '/wd/hub' },
        command: 'appium',
      },
    ],
  ],
  port: 4723,
  path: '/wd/hub',
  framework: 'mocha',
  reporters: ['spec'],
  mochaOpts: {
    ui: 'bdd',
    timeout: 120_000,
    retries: 1,
  },
};
