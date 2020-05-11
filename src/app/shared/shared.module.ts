import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

import { CardSmallComponent } from './components/card/small/small.component';
import { SectionContainerComponent } from './components/section/container/container.component';
import { SectionFooterComponent } from './components/section/footer/footer.component';
import { SectionHeaderComponent } from './components/section/header/header.component';

@NgModule({
  declarations: [
    CardSmallComponent,
    SectionContainerComponent,
    SectionFooterComponent,
    SectionHeaderComponent,
  ],
  imports: [CommonModule, RouterModule],
  exports: [
    CardSmallComponent,
    SectionContainerComponent,
    SectionFooterComponent,
    SectionHeaderComponent,
  ],
})
export class SharedModule {}
