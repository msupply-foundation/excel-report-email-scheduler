import React, { Component } from 'react';
import intl from 'react-intl-universal';

import { AppPluginMeta, PluginConfigPageProps } from '@grafana/data';
import { locales } from '../common/translations';

interface Props extends PluginConfigPageProps<AppPluginMeta> {}
type State = {
  shouldLoad: boolean;
};

/**
 * Factory creating class components which implement the interface
 * required by Grafana which requires a class, rather than functional
 * components (and no one likes classes!).
 */
export const ConfigPageFactory = (Content: any) =>
  class extends Component<Props, State> {
    displayName = 'ConfigPageFactory';
    constructor(props: Props) {
      super(props);

      this.state = {
        shouldLoad: false,
      };
    }

    async componentDidMount() {
      // TODO: More comprehensive localization solution
      await intl.init({ currentLocale: 'en', locales });
      this.setState({ shouldLoad: true });
    }

    render() {
      const { shouldLoad } = this.state;

      if (!shouldLoad) {
        return null;
      }

      return <Content {...this.props} />;
    }
  };
