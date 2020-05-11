import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { NavigationHeaderComponent } from './navigation/header/header.component';

@NgModule({
  declarations: [NavigationHeaderComponent],
  imports: [CommonModule, RouterModule],
  exports: [NavigationHeaderComponent],
})
export class CoreModule {}
