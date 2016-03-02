import { bindActionCreators } from 'redux'
import { createAction, Action } from 'redux-actions';
import { assign } from 'lodash';
import { Promise } from 'es6-promise';

import { File } from '../models/files';
import * as types from '../constants/ActionTypes';
import { socket } from './socket'

socket.on('ModifiedFile', (data) => {
  // dispatch(modifiedFile)
})

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
    console.log(socket)
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
    return socket.call('WatchFile', file)
      .then(file => dispatch(didWatchFile))
      .catch(err => console.error(err))
  }
}

export const unwatchFile = (file) => {
  return (dispatch, getState) => {
    return socket.call('UnwatchFile', file)
      .then(file => dispatch(didUnwatchFile))
      .catch(err => console.error(err))
  }
}

const changeFile = createAction<string>(
  types.CHANGE_FILE
)
