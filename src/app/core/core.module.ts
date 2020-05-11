import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';

import { NavigationHeaderComponent } from './navigation/header/header.component';

@NgModule({
  declarations: [NavigationHeaderComponent],
  imports: [CommonModule],
  exports: [NavigationHeaderComponent],
})
export class CoreModule {}
