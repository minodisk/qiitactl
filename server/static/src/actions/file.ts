import { createAction } from 'redux-actions';

import { File } from '../models/file'
import * as types from '../constants/ActionTypes';
import { Socket } from '../net/socket'

const socket:Socket = Socket.getInstance()

const willGetFile = createAction<void>(
  types.WILL_GET_FILE
)

const didGetFile = createAction<File>(
  types.DID_GET_FILE,
  (file:File) => file
);


export const getFile = (path: string) => {
  return (dispatch, getState) => {
    dispatch(willGetFile())
    return socket.call('GetFile', path)
      .then((file:File) => dispatch(didGetFile(file)))
      .catch(err => console.error(err))
  }
}

const didWatchFile = createAction<string>(
  types.DID_WATCH_FILE,
  (file:string) => file
);

const didUnwatchFile = createAction<string>(
  types.DID_UNWATCH_FILE,
  (file:string) => file
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
