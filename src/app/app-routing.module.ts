import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';

const routes: Routes = [
  {
    path: '',
    loadChildren: () =>
      import('./modules/home/home.module').then(
        (module: any) => module.HomeModule
      ),
  },
  {
    path: 'explore',
    loadChildren: () =>
      import('./modules/explore/explore.module').then(
        (module: any) => module.ExploreModule
      ),
  },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule],
})
export class AppRoutingModule {}
