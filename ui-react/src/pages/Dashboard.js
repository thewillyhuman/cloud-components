import React from 'react';
import {
  AppLayout, SideNavigation, ContentLayout,
  Header, Cards
} from '@cloudscape-design/components';

export default function Dashboard() {
  const services = [
    { id: 'ec2', name: 'EC2', description: 'Manage compute resources' },
    { id: 's3', name: 'S3', description: 'Scalable storage in the cloud' },
    { id: 'lambda', name: 'Lambda', description: 'Serverless compute functions' },
    { id: 'cloudwatch', name: 'CloudWatch', description: 'Monitor AWS resources and applications' },
  ];

  return (
    <AppLayout
      navigation={<SideNavigation items={[{ type: 'link', text: 'Dashboard', href: '/dashboard' }, { type: 'link', text: 'Dashboard', href: '/dashboard' }]} />}
      content={
        <ContentLayout header={<Header variant="h1">Available Services</Header>}>
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
      navigationOpen
    />
  );
}
