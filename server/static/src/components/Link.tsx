/// <reference path="../../typings/tsd.d.ts" />

import * as React from 'react'

interface LinkProps {
  path: string;
  active: boolean;
  actions: any;
}

interface LinkState {}

class Link extends React.Component<LinkProps, LinkState> {
  render() {
    const {path, active} = this.props
    if (active) {
      return <span>{path}</span>
    }
    return (
      <a
        href="#"
        onClick={(e) => {
          e.preventDefault()
        }}
      >
        {path}
      </a>
    )
  }
}

export default Link
