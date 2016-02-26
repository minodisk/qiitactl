import { combineReducers } from 'redux';
import { routerStateReducer } from 'redux-router'

import todos from './todos';
import file from './file';

const rootReducer = combineReducers({
  router: routerStateReducer,
  todos: todos,
  file: file
});

export { rootReducer };
