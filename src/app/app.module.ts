import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { CardSmallComponent } from './card-small/card-small.component';
import { HomeComponent } from './home/home.component';
import { NavigationBarComponent } from './navigation-bar/navigation-bar.component';
import { NewComponent } from './new/new.component';
import { SectionComponent } from './section/section.component';
import { SectionFooterComponent } from './section-footer/section-footer.component';
import { SectionHeaderComponent } from './section-header/section-header.component';
import { TrendingComponent } from './trending/trending.component';

@NgModule({
  declarations: [
    AppComponent,
    CardSmallComponent,
    HomeComponent,
    NavigationBarComponent,
    NewComponent,
    SectionComponent,
    SectionFooterComponent,
    SectionHeaderComponent,
    TrendingComponent,
  ],
  imports: [BrowserModule, AppRoutingModule],
  providers: [],
  bootstrap: [AppComponent],
})
export class AppModule {}
