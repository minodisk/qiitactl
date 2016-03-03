import { createAction } from 'redux-actions';

import { Tree } from '../models/tree';
import * as types from '../constants/ActionTypes';
import { Socket } from '../net/socket'

const socket:Socket = Socket.getInstance()

const willGetTree = createAction<void>(
  types.WILL_GET_TREE
)

const didGetTree = createAction<Tree>(
  types.DID_GET_TREE,
  (tree:Tree) => tree
);

export const getTree = () => {
  return (dispatch, getState) => {
    dispatch(willGetTree())
    return socket.call('GetTree', null)
      .then(tree => dispatch(didGetTree(tree)))
      .catch(err => console.error(err))
  }
}
