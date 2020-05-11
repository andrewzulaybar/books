import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

import { NewComponent } from './pages/new/new.component';
import { TrendingComponent } from './pages/trending/trending.component';

const routes: Routes = [
  {
    path: 'new',
    component: NewComponent,
  },
  {
    path: 'trending',
    component: TrendingComponent,
  },
];

@NgModule({
  imports: [RouterModule.forChild(routes)],
  exports: [RouterModule],
})
export class ExploreRoutingModule {}
