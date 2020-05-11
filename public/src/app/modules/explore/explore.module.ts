import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { ExploreRoutingModule } from './explore-routing.module';

import { NewComponent } from './pages/new/new.component';
import { TrendingComponent } from './pages/trending/trending.component';

@NgModule({
  declarations: [NewComponent, TrendingComponent],
  imports: [CommonModule, ExploreRoutingModule],
  exports: [NewComponent, TrendingComponent],
})
export class ExploreModule {}
