import { createAction, Action } from 'redux-actions';
import { assign } from 'lodash';
import { Promise } from 'es6-promise';

import { File } from '../models/files';
import * as types from '../constants/ActionTypes';
import { socket } from './socket'


const requestFile = createAction<void>(
  types.REQUEST_FILES
)

const recieveFile = createAction<File>(
  types.RECIEVE_FILES,
  (file: File) => file
);

export const fetchFiles = () => {
  return (dispatch, getState) => {
    dispatch(requestFile())
    return socket.call('GetAllFiles', null)
      .then(file => dispatch(recieveFile(file)))
  }
}

const didWatchFile = createAction<string>(
  types.DID_WATCH_FILE,
  (file: string) => file
);

const didUnwatchFile = createAction<string>(
  types.DID_UNWATCH_FILE,
  (file: string) => file
);

export const watchFile = (file) => {
  return (dispatch, getState) => {
    console.log('watchFile:', file)
    return socket.call('WatchFile', file)
      .then(file => dispatch(didWatchFile))
  }
}

export const unwatchFile = (file) => {
  return (dispatch, getState) => {
    console.log('unwatchFile:', file)
    return socket.call('UnwatchFile', file)
      .then(file => dispatch(didUnwatchFile))
  }
}
