/// <reference path="../../typings/tsd.d.ts" />

import { combineReducers } from 'redux';

import paths from './paths';

const rootReducer = combineReducers({
  paths: paths
});

export { rootReducer };
