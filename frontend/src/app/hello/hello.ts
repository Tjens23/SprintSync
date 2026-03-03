import { Component } from '@angular/core';

@Component({
  selector: 'app-hello',
  imports: [],
  templateUrl: './hello.html',
  styleUrls: ['./hello.css'],
})
export class Hello {
  hello: string = 'Hello world from Angular!';
  isVisible: boolean = true;

  toggleVisibility() {
    this.isVisible = !this.isVisible;
  }
}
