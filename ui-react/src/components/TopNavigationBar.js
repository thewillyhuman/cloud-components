import React from 'react';
import { TopNavigation } from '@cloudscape-design/components';

export default function TopNavigationBar() {
  return (
    <TopNavigation
      identity={{ href: '/', title: 'Cloud Manager' }}
      utilities={[]}
    />
  );
}
