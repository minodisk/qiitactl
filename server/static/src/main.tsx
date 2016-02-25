import * as React from 'react';
import * as ReactDOM from 'react-dom';

import {
  Store,
  compose,
  createStore,
  bindActionCreators,
  combineReducers,
  applyMiddleware
} from 'redux';
import {
  connect,
  Provider
} from 'react-redux';
import * as thunkMiddleware from 'redux-thunk'
const createLogger = require('redux-logger')

import { rootReducer } from './reducers/rootReducer';
import { openSocket  } from './actions/socket'
import { fetchFiles } from './actions/files'
import App from './containers/App';

const loggerMiddleware = createLogger()
const store: Store = createStore(rootReducer,
  applyMiddleware(
    thunkMiddleware,
    loggerMiddleware
  )
);

store
  .dispatch(openSocket())
  .then(() => store.dispatch(fetchFiles()))
  .then(() => console.log('done:', store.getState()))

ReactDOM.render(
  <Provider store={store}>
    <App/>
  </Provider>,
  document.getElementById('app')
);
