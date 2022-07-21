import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FontAwesomeModule } from '@fortawesome/angular-fontawesome';

import { HomeComponent } from './home/home.component';
import { UsersComponent } from './users/users.component';

@NgModule({
  imports: [CommonModule, FontAwesomeModule],
  declarations: [HomeComponent, UsersComponent],
  exports: [HomeComponent, UsersComponent],
  providers: [],
})
export class ScreensModule {}
