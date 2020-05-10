import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { HomeComponent } from './home/home.component';
import { NewComponent } from './new/new.component';
import { TrendingComponent } from './trending/trending.component';

const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'new', component: NewComponent },
  { path: 'trending', component: TrendingComponent },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
