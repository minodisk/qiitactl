import { combineReducers } from 'redux';
import { routerStateReducer } from 'redux-router'

import todos from './todos';
import tree from './tree';

const rootReducer = combineReducers({
  router: routerStateReducer,
  todos: todos,
  tree: tree,
});

export { rootReducer };
