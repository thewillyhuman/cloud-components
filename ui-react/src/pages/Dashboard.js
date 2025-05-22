import React from 'react';
import {
  AppLayout, ContentLayout,
  Header, Cards
} from '@cloudscape-design/components';
import SideNavigation, { SideNavigationProps } from '@cloudscape-design/components/side-navigation';
import StatusIndicator from "@cloudscape-design/components/status-indicator";

const navItems: SideNavigationProps.Item[] = [
    {
      type: 'section',
      text: 'Cell',
      items: [
        { type: 'link', text: 'Configuration', href: '#/pages' },
        { type: 'link', text: 'Controllers', href: '#/pages' },
        { type: 'link', text: 'Workers', href: '#/users' },
      ],
    },
    {
      type: 'section',
      text: 'Hardware',
      items: [
        { type: 'link', text: 'Servers', href: '#/pages' },
        { type: 'link', text: 'Disks', href: '#/pages' },
        { type: 'link', text: 'NICs', href: '#/pages' },
        { type: 'link', text: 'Switches', href: '#/users' },
      ],
    },
    {
      type: 'section',
      text: 'Network',
      items: [
        { type: 'link', text: 'DHCP', href: '#/database' },
        { type: 'link', text: 'DNS', href: '#/authentication' },
        { type: 'link', text: 'Route Tables', href: '#/authentication' },
        { type: 'link', text: 'Load Balancers', href: '#/authentication' },
      ],
    },
    {
      type: 'section',
      text: 'Storage',
      items: [
        { type: 'link', text: 'Volumes', href: '#/authentication' },
        { type: 'link', text: 'File Systems', href: '#/authentication' },
      ],
    },
    {
      type: 'section',
      text: 'Compute',
      items: [
        { type: 'link', text: 'Jobs', href: '#/database' },
        { type: 'link', text: 'Images', href: '#/database' },
      ],
    },
    {
      type: 'section',
      text: 'Databases',
      items: [
        { type: 'link', text: 'Cockroach', href: '#/database' },
        { type: 'link', text: 'Cassandra', href: '#/database' },
        { type: 'link', text: 'Mongo', href: '#/database' },
        { type: 'link', text: 'ElasticSearch', href: '#/database' },
      ],
    },
    {
      type: 'section',
      text: 'Monitoring',
      items: [
        { type: 'link', text: 'Mimir', href: '#/database' },
        { type: 'link', text: 'Loki', href: '#/database' },
        { type: 'link', text: 'Tempo', href: '#/database' },
      ],
    },
  ];

export default function Dashboard() {
  const services = [
    { id: 'ec2', name: 'EC2', description: 'Manage compute resources' },
    { id: 's3', name: 'S3', description: 'Scalable storage in the cloud' },
    { id: 'lambda', name: 'Lambda', description: 'Serverless compute functions' },
    { id: 'cloudwatch', name: 'CloudWatch', description: 'Monitor AWS resources and applications' },
  ];

  return (
    <AppLayout
      navigation={<SideNavigation items={navItems} />}
      content={
        <ContentLayout header={<Header variant="h1">Available Services <StatusIndicator type="error">Error</StatusIndicator></Header>}>
          <Cards
            cardsPerRow={[{ cards: 2 }]}
            items={services}
            cardDefinition={{
              header: item => item.name,
              sections: [{ id: 'description', content: item => item.description }]
            }}
          />
        </ContentLayout>
      }
      toolsHide={true}
      
    />
  );
}
