import * as React from 'react';

import TreeElement from './TreeElement'
import * as tree from '../models/tree'

const styles = require('../styles/file.css')

interface Props {
  indent:number;
  tree:tree.Tree;
}

interface State {
  opened: boolean;
}

export default class Tree extends React.Component<Props, State> {
  constructor(props, context) {
    super(props, context)
    this.state = {
      opened: true,
    }
  }

  toggleOpen() {
    const {opened} = this.state
    this.setState({
      opened: !opened,
    })
  }

  render() {
    const {indent, tree} = this.props
    if (tree.children == null) {
      return (
        <div>
          <TreeElement
            indent={indent}
            className={styles.file}
            active={true}
            title={tree.abs}
            children={tree.name}
            linkTo={tree.rel}
          />
        </div>
      )
    }
    return (
      <div>
        <TreeElement
          indent={indent}
          className={this.state.opened ? styles.dirOpened : styles.dirClosed}
          active={true}
          title={tree.abs}
          children={tree.name}
          onClick={() => this.toggleOpen()}
        />
        <ul
          className={this.state.opened ? styles.filesOpened : styles.filesClosed}
        >
          {
            tree.children.map((t:tree.Tree) => (
              <li key={t.id}>
                <Tree indent={indent + 1} tree={t} />
              </li>
            ))
          }
        </ul>
      </div>
    )
  }
}
